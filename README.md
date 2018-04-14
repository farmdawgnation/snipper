# Snipper

Snipper is a proof of concept for applying Lift Snippet semantics to YAML
templating.

Snipper takes in a template and then any number of transformers to create a
final YAML file. YAML transformers are other YAML files whose keys conform to
selector syntax and values indicate what input to the provided selector should
be used.

## Installing

Snipper can be installed by downloading a binary for your system from the releases
page. Once installed somewhere useful on your path I recommend checking out the usage
instructions first:

```
snipper - snippet style transformers for YAML
Usage: snipper template.yaml transformer.yaml [transformer.yaml [transformer.yaml ...]]
```

## Using

Snipper takes in at least two yaml files: a template and one or more transformers that are
applied to the template in the order they are provides.

Transformers are YAML files with special formatting rules. Those rules are:

* The Transformer YAML structure must not nest: it must be a single-level object with keys
  and values.
* The keys of the Transformer YAML are specially encoded to represent nested items.

If you want to see what I mean, check out the examples folder where I've outlined a few different
ways the Transformer YAML works.

## Limitations

This is very much a proof of concept at this point. If you like the idea, please star it. Of note
snipper doesn't really have any support for:

* Altering objects nested in arrays in the template
* Probably many other things I'm not thinking of this moment?

## Building

To build this project you'll need:

* Go (1.9 or greater)
* Dep

Once cloned run `dep ensure` to get the dependencies. Then `go build` to your heart's content.
I'm not a fan of vendored repositories, so please restrain yourself from opening a PR or issue
aong the lines of "you should have checked in your vendor directory."

## About the author

Matt Farmer works at [MailChimp](https://mailchimp.com) and occasionally blogs at
[farmdawgnation.com](https://farmdawgnation.com).
