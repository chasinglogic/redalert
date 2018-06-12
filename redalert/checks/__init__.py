from .exc import CheckFailure
from .file_exists import FileExistsCheck, FileNotExistsCheck
from .ulimit_checks import AddressSizeCheck, UlimitCheck, OpenFilesCheck


def get_check(name, args=None):
    '''Return the appropriate check instance based on test name.'''
    if args is None:
        args = {}

    if name == 'address-size':
        return AddressSizeCheck(**args)
    elif name == 'open-files':
        return OpenFilesCheck(**args)
    elif name == 'file-exists':
        return FileExistsCheck(**args)
    elif name == 'file-does-not-exist':
        return FileNotExistsCheck(**args)
    elif name == 'dpkg-installed':
        from .dpkg_installed import DpkgCheck
        return DpkgCheck(**args)
    elif name == 'compile-gcc':
        from .compile_gcc import CompileGccCheck
        return CompileGccCheck(**args)

    raise NotImplementedError
