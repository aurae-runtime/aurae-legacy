/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package core

// CoreServicer is the interface that defines the basic methods that all
// CoreService implementations must use.
//
// This is a notable offensive interface as it (by design) has no concept of
// error handling.
//
// The error-less convention exists to provide simple guarantees to the
// calling code that will guarantee that a database error will not crash the
// calling code.
//
// For atomic data transactions, the best-effort approach will be to call Set()
// followed quickly by a Get() to ensure the data exists as it was set.
//
// Other than this "best-effort" approach there are no persistence guarantees
// with any of the CoreServicer implementations. At least not through this interface.
type CoreServicer interface {

	// Get will return a value for a key (if it exists)
	Get(key string) string

	// Set will set a key to a value
	Set(key, value string)

	// List is a map[filename]isFile where isFile=true if the dirent is a file, or false
	// for a directory.
	List(key string) map[string]bool

	// Remove is a recursive-by-default function that will remove a given node (unless root)
	// and all of its children.
	Remove(key string)
}
