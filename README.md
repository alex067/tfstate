# Tfstate

A wrapper around terraform state commands.

## Description 

tfstate is a wrapper around terraform state commands which alter the state file, such as ```terraform state mv``` and ```terraform state rm```

tfstate automatically generates a backup inside ```.terraform/tfstate/*``` allowing for easy rollback if a mistake was made. This is possible due to the wrapper, by first performming a backup of the current state file, then running ```tfstate rollback --latest```

tfstate adds a manual confirmation step, which also lists the possible resources affected by the requested state command.

Example:
```hcl

resource "null_resource" "main" {}
resource "null_resource" "deps" {
  depends_on = [null_resource.main]
}
```

Running ```tfstate rm null_resource.main``` generates an ouptut containing a list of resources affected by the command, such as:
```
null_resource.deps
... Affected resources: 1 
```

## Usage

tfstate is a small wrapper around terraform state. To run the tool, simply use ```tfstate``` instead of ```terraform state``` for the following commands:
* ```terraform state mv``` -> ```tfstate mv```
* ```terraform state rm``` -> ```tfstate rm```

You may also create a backup of your state file
```tfstate backup```

To rollback to a specific state file, you must first have ran your ```terraform state rm/mv``` using tfstate, as tfstate generates backups in ```.terraform/tfstate```

Select a state file in ```.terraform/tfstate``` to rollback to
```tfstate rollback state-0-123456.json```

## Download 

You may download the latest binary by visitng the tags and downloading the artifact. The zip file contains binary builds for Windows and Darwin machine types. Other architectures such as linux has not been tested yet but I'll get to it soon!
