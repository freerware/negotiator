---
layout: post
title:  "Using Protobuf With negotiator"
categories: posts
---
Among the many ways to transmit representations, Google's
[Protobuf][protobuf-docs] (short for "protocol buffer") has gained traction.
With Protobuf, your representations (referred to as "messages" in the Protobuf 
documentation) are defined in `.proto` files. Using the Protobuf compiler, the 
definitions in the `.proto` files are then generated in your language(s) of
choice. 

> For a quick tutorial on using Protobuf with Go, check out
[this tutorial][protobuf-go-tutorial].

So can we use Protobuf with `negotiator`? You bet! Let's walk through it.

## Create your `.proto` file

First things first - let's create our `.proto` file containing the
definitions of the Protobuf messages. Suppose we have
an `account` resource for our API with this definition:

{% highlight proto %}
// you need to use Protobuf version 3 in order to support Go.
syntax = "proto3";
// the package for our Protobuf messages.
package tutor;
// import other message definitions, such as for timestamps.
import "google/protobuf/timestamp.proto";
// the full import path of the Go package that contains the generated code.
option go_package = "github.com/freerware/tutor/api/representations/protobuf/gen";

message Account {
  string UUID                         = 1;
  string username                     = 2;
  string givenName                    = 3;
  string surname                      = 4;
  google.protobuf.Timestamp createdAt = 5;
  google.protobuf.Timestamp updatedAt = 6;
  google.protobuf.Timestamp deletedAt = 7;
}
{% endhighlight %}

## Generate your Go types

Once you have [downloaded the Protobuf compiler][protobuf-compiler-download],
we need to install the Go protobuf plugin:

{% highlight bash %}
go install google.golang.org/protobuf/cmd/protoc-gen-go
{% endhighlight %}

Next, we need to invoke the compiler to generate the corresponding Go types
that represent our messages. Based on our project structure, we invoked the
compiler like this:

{% highlight bash %}
protoc --proto_path=$PROJ_DIR --go_opt=paths=source_relative --go_out=$PROJ_DIR $PROJ_DIR/api/representations/protobuf/gen/tutor.proto
{% endhighlight %}

The `$PROJ_DIR` is the directory for our project. The `--proto_path` option
specifies a directory in which to search for imports. The `--go_opt=paths=source_relative`
option allows us to place our generated source code in the same directory as
our `.proto` file. `--go_out` specifies the base directory in which to place
the generated source code. Finally, the final piece of the command is an
argument specifying the path the `.proto` file. 

> For more details, head over to the [official documentation][protoc-go-docs] on how to use 
the `protoc` compiler for Go.

## Create your representations

Now that we have generated our Protobuf messages, we need to define our representations.
We do this by creating a Go type that implements
[`representation.Representation`][representation-docs], and embeds our Protobuf
generated type as well as [`representation.Base`][representation-base-docs]:

{% highlight golang %}
package protobuf

import (
	"errors"

	"github.com/freerware/negotiator/representation"
	"github.com/freerware/tutor/api/representations/protobuf/gen"
	"github.com/freerware/tutor/domain"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// The media types representing Protobuf content.
const (
	mediaTypeProtobuf  = "application/protobuf"
	mediaTypeXProtobuf = "application/x-protobuf"
)

// Account represents the Protobuf representation for an account resource.
type Account struct {
	representation.Base
	gen.Account // the generated Protobuf type.
}

// NewAccount constructs a new account representation.
func NewAccount(a domain.Account) Account {
	// define a customer marshaller.
	marshaller := func(in interface{}) ([]byte, error) {
		message, ok := in.(proto.Message)
		if !ok {
			return []byte{}, errors.New("must provide Protobuf message to marshal successfully")
		}
		return proto.Marshal(message)
	}

	// define a custom unmarshaller.
	unmarshaller := func(b []byte, out interface{}) error {
		message, ok := out.(proto.Message)
		if !ok {
			return errors.New("must provide Protobuf message to unmarshal successfully")
		}
		return proto.Unmarshal(b, message)
	}

	// set the state of the representation based on the provided entity.
	acc := Account{}
	acc.UUID = a.UUID().String()
	acc.GivenName = a.GivenName()
	acc.Surname = a.Surname()
	acc.Username = a.Username()
	c := timestamppb.Timestamp{Seconds: a.CreatedAt().Unix()}
	acc.CreatedAt = &c
	u := timestamppb.Timestamp{Seconds: a.UpdatedAt().Unix()}
	acc.UpdatedAt = &u
	var d *timestamppb.Timestamp
	if a.DeletedAt() != nil {
		d = &timestamppb.Timestamp{Seconds: a.DeletedAt().Unix()}
	}
	acc.DeletedAt = d

	// set representation metadata.
	acc.SetContentCharset("ascii")
	acc.SetContentLanguage("en-US")
	acc.SetContentType(mediaTypeProtobuf)
	acc.SetSourceQuality(1.0)
	acc.SetContentEncoding([]string{"identity"})

	// set the custom marshaller.
	acc.SetMarshallers(map[string]representation.Marshaller{
		mediaTypeProtobuf:  marshaller,
		mediaTypeXProtobuf: marshaller,
	})

	// set the custom unmarshaller.
	acc.SetUnmarshallers(map[string]representation.Unmarshaller{
		mediaTypeProtobuf:  unmarshaller,
		mediaTypeXProtobuf: unmarshaller,
	})
	return acc
}

// Bytes serializes the representation.
func (a Account) Bytes() ([]byte, error) {
	return a.Base.Bytes(&a)
}

// FromBytes deserializes the representation.
func (a *Account) FromBytes(b []byte) error {
	return a.Base.FromBytes(b, a)
}
{% endhighlight %}

## Negotiate!

That's it! With our HTTP handler below, we now have everything we need to
negotiate using Protobuf:

{% highlight golang %}
//...

func (ar *AccountResource) Get(w http.ResponseWriter, request *http.Request) {

	// retrieve the account uuid.
	vars := mux.Vars(request)
	uuid, err := u.FromString(vars["uuid"])
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	// retrieve the account.
	account, err := ar.accountService.Get(uuid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jacc := j.NewAccount(account)
	jacc.SetContentLocation(*request.URL)
	gjacc := j.NewAccount(account)
	gjacc.SetContentLocation(*request.URL)
	gjacc.SetContentEncoding([]string{"gzip"})
	yacc := y.NewAccount(account)
	yacc.SetContentLocation(*request.URL)
	xacc := x.NewAccount(account)
	xacc.SetContentLocation(*request.URL)
	pacc := p.NewAccount(account)
	pacc.SetContentLocation(*request.URL)
	representations := []representation.Representation{jacc, yacc, xacc, gjacc, pacc}

	// negotiate.
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: w}
	if err = proactive.Default.Negotiate(ctx, representations...); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

//...
{% endhighlight %}

## Finished code

Much of the code in this post was lifted from [`tutor`][tutor], our sample RESTful API
demonstrating example use of the `freerware` product suite.

[protobuf-docs]: https://developers.google.com/protocol-buffers
[protobuf-go-tutorial]: https://developers.google.com/protocol-buffers/docs/gotutorial
[protobuf-compiler-download]: https://developers.google.com/protocol-buffers/docs/downloads
[protoc-go-docs]: https://developers.google.com/protocol-buffers/docs/reference/go-generated
[tutor]: https://github.com/freerware/tutor
[representation-docs]: https://github.com/freerware/negotiator/blob/master/representation/representation.go
[representation-base-docs]: https://github.com/freerware/negotiator/blob/master/representation/base.go
