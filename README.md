# Snipper

Snipper is a project that started with the question "What if managing a bunch
of YAML for sophisticated Helm Charts didn't have to involve Go text templating?"

Inspired by the [Lift Framework][lift] snippets, Snipper aims to separate the
**structure file** from the **transformations** that need to be applied to it.
The most important thing: _transformations are YAML too._

Snipper supports:

* Explicitly setting a particular key in an object from your structure.
* Explicitly setting a particular member of an array from your structure.
* Making changes to all members of an array from your structure.
* Making changes to arrays of objects that have a certain property set to a
  specified value. (e.g. only change `image=alpine` containers)
* Appending text to string values.
* Appending members to arrays.

Snipper does not support:

* **If statements.** If you have transforms that should only be applied in some
  situations, separate them out into their own file and do or don't apply them
  to your template. Admittedly, this moves the "conditional" aspect of rendering
  your YAML to somewhere else. However, in exchange you gain a significant boost
  in long term readability if your YAML grows and gains functionality and
  considerations, as most things do. Further, separating out the transforms into
  separate projects later on is a lot easier than disentangling a bunch of `if`/
  `elseif` and ensuring that you got the same result.

[lift]: http://liftweb.net

## Installing

Snipper can be installed by downloading a stable binary for your system from the
[releases page][releases]. Once installed somewhere useful on your `PATH` I
recommend checking out the usage instructions first:

```
snipper - snippet style transformers for YAML
Usage: snipper template.yaml transformer.yaml [transformer.yaml [transformer.yaml ...]]
```

Alternately, if you want to work with the cutting edge code and you have a
working Go environment, you can do the following:

```
$ go get github.com/farmdawgnation/snipper
$ go install github.com/farmdawgnaton/snipper
```

[releases]: https://github.com/farmdawgnation/snipper/releases

## Using Snipper

Using Snipper from the command line is designed to be pretty straighforward.
Simply invoke the binary, provide a template YAML file, and then any number of
transformers that will be applied in that order.

Transformers are just YAML files, too, with a minor twist: the top level keys
of the YAML object are interpreted as instructions for Snipper. They are called
**selectors.** Conceptually, these are very similar to the CSS Selectors that
have been in use on the web for a long time. They merely have a tweaked syntax
that makes more sense in YAML-land. The selectors describe the "coordinates" of
where the **value** (the thing that ends up in the rendered YAML) should go.

Let's take a concrete example. I've got some YAML describing a dog named Shadow.

```yaml
# shadow-template.yaml
name: Shadow
type: Doggo
```

Shadow is a very good dog, and I want her YAML to indicate that. I could apply
the following selector to add a `goodDog` key to Shadow's YAML:

```yaml
# good-dog.yaml
'goodDog': true
```

If I applied this to my YAML I would end up with:

```yaml
goodDog: true
name: Shadow
type: Doggo
```

We could even get a bit more complex and replace her `type` with an array using
this transform:

```yaml
# collie-dogs.yaml
'type':
- Doggo
- Collie
```

And I could invoke snipper and ask it to apply both of these to my template:

```
$ snipper shadow-template.yaml good-dog.yaml collie-dogs.yaml
goodDog: true
name: Shadow
type:
- Doggo
- Collie
```

Snipper successfully separates the rules that need to be applied to a template
from the template itself.

## Selector Syntax in Detail

Snipper's selector syntax has the following rules:

* Selectors will almost always need to be single-quoted.
* By default, selectors are selecting on keys of YAML objects. In the example
  above, just providing `'goodDog'` as the selector was good enough to get the
  effect we wanted.
* You can select keys nested in other objects by delimiting those with a colon
  (`:`). In a Kubernetes pod spec the containers are at `'spec:containers'` in
  Snipper selector syntax. We chose this because periods (`.`) are exceptionally
  common field names in K8S-land. The piece on either side of a colon is called
  a **path**.
* If a path has square brackets in it (`[]`) it does something with arrays.
  * A path of `[]` just selects all arrays at a location. So, `spec:containers:[]`
    hits all of the containers in your pod.
  * An array path with a number in it (such as `[0]`) selects a specific element
    of the array.
  * An array paths can also only match array members with a certain property.
    For example, `spec:containers:[image=alpine]` only matches members of your
    containers array that have the value of `image` set to `alpine`.

Selectors will, by default, replace the value at the coordinates they select
with the value that appears in the transform definition. You can alter this
behavior by adding an **action character** to the end of your selector:

* Currently, the only supported action character is `+` which causes the
  selector to append its value to the location in the template. So if you
  wanted to add an explicit tag to all your alpine containers you could use
  this: `'spec:containers:[image=alpine]:image+': latest`.

We encourage you to browse the [examples][examples] folder and try some of them
out for yourself to see what they do.

[examples]:https://github.com/farmdawgnation/snipper/tree/master/examples

## Building

To build this project you'll need:

* Go (1.9 or greater)
* Dep

Once cloned run `dep ensure` to get the dependencies. Then `go build` to your
heart's content. I'm not a fan of vendored repositories, so please restrain
yourself from opening a PR or issue along the lines of "you should have checked
in your vendor directory."

## About the author

Matt Farmer works at [MailChimp](https://mailchimp.com) and occasionally blogs at
[farmdawgnation.com](https://farmdawgnation.com).
