# Obsidify
Converts Quiver App markdown to Obsidian MD markdown.

This was built out of immediate frustration when attempting to try Obsidian, and I figure it may save someone else some time.

To use this, you need to export your notes in Quiver by grouping into directories in markdown format (right click > export > markdown). 

If you export your notes all at once, you'll lose your structure (which you can do if you want that). Otherwise, you need to export by grouping. Quiver will create the directories for you when exporting.

The ideal method of solving this would be to parse the `qvnote` files. This was designed to get the job done with minimal effort.

Before running this script, in Settings under Editor set the following:
`Default location for new attachments` to `In subfolder under current folder`
`Subfolder name` to `attachments`

This will do the following:
- Fix embedded image links
- Fix embedded document links and rename the file from the Quiver generated hex string to the name in the caption you used
- Fix internal links so they continue to work
- Rename resources to `attachments`
- Remove any duplicates (notes with the same name, in a parent directory). You can comment this code out as needed, I noted the relevant section in `main.go`

Before running this put all your exported Quiver App notes into `./source`. Obsidify will update the files and directories within `./source`.

Build (`go build`) or run (`go run main.go`). Example output:

```
$ go run main.go
* processing
 + source/Swift 3/Print headers, request body, response body.md
 + source/Swift 3/Programmatically calling a segue, and send data with it.md
    ![IMAGE](resources/5BC8900B4FD4E470AD0F86581F4DAF7C.jpg) -> ![[attachments/5BC8900B4FD4E470AD0F86581F4DAF7C.jpg]]
    ![Layer.tiff](resources/0A7BC4066F273465AD9C74EB041030B8.jpg) -> ![[attachments/0A7BC4066F273465AD9C74EB041030B8.jpg]]
 + source/Swift 3/Serial queue example.md
 + source/Swift 3/Set AppIcon badge number.md
 + source/Swift 3/Set navigation bar style globally.md
 + source/Swift 3/Simple NSOperations example.md
 + source/Swift 3/Simple semaphore example.md
 + source/XCode/Versioning (CFBundleVersion).md
 ...
* renaming resource directories
* cleaning up duplicates based on filenames and length of paths
 - source/Test/test.md
 - source/Security/recon-ng.md
 ...
* notes converted: 499, duplicates deleted: 210, actual note count: 289
```