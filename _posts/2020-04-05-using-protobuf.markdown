---
layout: post
title:  "Using Protobuf With negotiator"
categories: posts
---
Among the many ways to transmit representations, Google's [Protobuf][protobuf-docs] (short for "protocol buffer") has gained traction. With Protobuf, your representations (referred to as "messages" in the Protobuf community) are defined in `.proto` files. Using the Protobuf compiler, the definitions in the `.proto` files are then generated in your language(s) of choice. These generated types represent the representation itself. For a quick tutorial on using Protobuf with Go, check out [this tutorial][protobuf-go-tutorial].

So can we use Protobuf with `negotiator`? You bet! Let's walk through it.

## Create your `.proto` file

First things first - let's create our `.proto` file containing the definitions of the representations ("messages"). Suppose we have an `account` resource for our API with this representation:

{% highlight proto %}
// you need to use Protobuf version 3 in order to support Go.
syntax = "proto3";
// the package for our Protobuf representations.
package tutor;

// import other representations (messages), such as for timestamps.
import "google/protobuf/timestamp.proto";

option go_package = "github.com/freerware/tutor/api/representations/protobuf/go/tutorpb";

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

Once you have [downloaded the Protobuf compiler][protobuf-compiler-download], we need to install the Go protobuf plugin:

{% highlight bash %}
go install google.golang.org/protobuf/cmd/protoc-gen-go
{% endhighlight %}

Finally, invoke the compiler to generate your Go types:

{% highlight bash %}
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/api/representations/protobuf/tutor.proto
{% endhighlight %}

The `$SRC_DIR` is the directory for your project, and `$DST_DIR` is where you would like to store the generated Go files.

## Create your representations

Now that we have generated your Protobuf messages, we need to define our representations. We do this by creating a Go type that implements [`representation.Representation`][representation-docs], and embeds your Protobuf generated type as well as [`representation.Base`][representation-base-docs]:

{% highlight golang %}
type Account struct {
  tutorpb.Account
  representation.Base
}

// NewAccount constructs a new account representation.
func NewAccount(a domain.Account) Account {
  marshaller := func(in interface{}) ([]byte, error) {
    message, ok := in.(proto.Message)
    if !ok {
      return []byte{}, errors.New("must provide Protobuf message to marshal successfully"
    }
    return proto.Marshal(message)
  }
  unmarshaller := func(b []byte, out interface{}) error {
    message, ok := out.(proto.Message)
    if !ok {
      return []byte{}, errors.New("must provide Protobuf message to unmarshal successfully"
    }
    return proto.Unmarshal(b, message)
  }
  a := Account{
    UUID:              a.UUID(),
    GivenName:         a.GivenName(),
    Surname:           a.Surname(),
    PrimaryCredential: a.Username(),
    CreatedAt:         a.CreatedAt(),
    UpdatedAt:         a.UpdatedAt(),
    DeletedAt:         a.DeletedAt(),
  }
  a.SetContentCharset("ascii")
  a.SetContentLanguage("en-US")
  a.SetContentType("application/protobuf")
  a.SetSourceQuality(1.0)
  a.SetContentEncoding([]string{"identity"})
  a.SetMarshallers(map[string]representation.Marshaller {
    "application/protobuf": marshaller,
    "application/x-protobuf": marshaller,
  })
  a.SetUnmarshallers(map[string]representation.Unmarshaller {
    "application/protobuf": unmarshaller,
    "application/x-protobuf": unmarshaller,
  })
  return a
}

func (a Account) Bytes() ([]byte, error) {
  return a.Base.Bytes(&a)
}

func (a Account) FromBytes(b []bytes) error {
  return a.Base.FromBytes(&a, b)
}
{% endhighlight %}

## Negotiate!

That's it! You now have everything you need to negotiate using Protobuf:

{% highlight golang %}
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
  
  // JSON.
  jacc := j.NewAccount(account)
  jacc.SetContentLocation(*request.URL)
  
  // JSON + gzip.
  gjacc := j.NewAccount(account)
  gjacc.SetContentLocation(*request.URL)
  gjacc.SetContentEncoding([]string{"gzip"})
  
  // YAML.
  yacc := y.NewAccount(account)
  yacc.SetContentLocation(*request.URL)
  
  // XML.
  xacc := x.NewAccount(account)
  xacc.SetContentLocation(*request.URL)
  
  // Protobuf.
  pacc := protobuf.NewAccount(account)
  pacc.SetContentLocation(*request.URL)
  
  representations := []representation.Representation{jacc, yacc, xacc, gjacc, pacc}
  
  // negotiate.
  ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: w}
  if err = proactive.Default.Negotiate(ctx, representations...); err != nil {
  	http.Error(w, err.Error(), 500)
  }
}
{% endhighlight %}

## Finished code

If you want to check out the code featured in this post in action, check out [`tutor`][tutor].

[protobuf-docs]: https://developers.google.com/protocol-buffers
[protobuf-go-tutorial]: https://developers.google.com/protocol-buffers/docs/gotutorial
[protobuf-compiler-download]: https://developers.google.com/protocol-buffers/docs/downloads
[tutor]: https://github.com/freerware/tutor
[representation-docs]: https://github.com/freerware/negotiator/blob/master/representation/representation.go
[representation-base-docs]: https://github.com/freerware/negotiator/blob/master/representation/base.go
