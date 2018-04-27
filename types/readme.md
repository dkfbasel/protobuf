# Initialize protobuf types

To use the protobuf types, they have to be linked to the ```/usr/local/include```
folder, in the same way like Google's protobuf well-known-types.

```console
foo@bar:~$ mkdir /usr/local/include/dkfbasel
foo@bar:~$ ln -s $GOPATH/src/github.com/dkfbasel/protobuf/types /usr/local/include/dkfbasel/protobuf
```

Now the protobuf types can be imported into your proto definition file.
```
syntax = "proto3";

package startrek;

import "dkfbasel/protobuf/timestamp.proto";
import "dkfbasel/protobuf/nullstring.proto";
import "dkfbasel/protobuf/nullint.proto";

message StarfleetShip {
	string name = 1;
	dkfbasel.protobuf.NullInt passengers = 2;
	dkfbasel.protobuf.NullString mission = 3;
	dkfbasel.protobuf.Timestamp departure_time = 4;
}
```
