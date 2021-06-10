# zoom-splitter

--- This is a work in progress, not functional yet! ---

My business does video interviews on Zoom. But Zoom recordings are split-screen, and we want them as separate videos for post-processing.

This project is a small container which receives a file path in an HTTP push, splits the targeted (mounted, external) file, and moves the results to a given (mounted, external) directory. 

The intended use case is that we can upload a video to our Nextcloud instance, NC Flow will hit the HTTP endpoint, and this container will do the rest.

This is a home hobby project and will be slowly developed. No warranty or support implied. 

