// Package folder provides utilities for working with file system directories.
// It defines a Folder type which is a string representing a directory path.
// The package includes methods for creating, removing, expanding, and checking
// the existence of directories, as well as manipulating paths.
//
// Example usage:
//
//	f := folder.New("/tmp", "example")
//	if err := f.Create(); err != nil {
//	    log.Fatal(err)
//	}
//	defer f.Remove()
//
//	fmt.Println("Folder created:", f.Path())
package folder
