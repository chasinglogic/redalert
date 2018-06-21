import pytest

from redalert.checks.exc import CheckFailure
from redalert.checks.python_module_version import PythonModuleCheck


@pytest.mark.agnostic
def test_python_module_version_eq():
    args = {'module': 'click', 'version': '6.7', 'comparison': 'eq'}
    check = PythonModuleCheck(**args)
    check.check()

    args['version'] = '0.0.0'

    check = PythonModuleCheck(**args)
    with pytest.raises(CheckFailure):
        check.check()


@pytest.mark.agnostic
def test_python_module_version_gte():
    args = {'module': 'click', 'version': '6.8', 'comparison': 'gte'}
    check = PythonModuleCheck(**args)
    with pytest.raises(CheckFailure):
        check.check()

    args['version'] = '0.0.0'

    check = PythonModuleCheck(**args)
    check.check()


@pytest.mark.agnostic
def test_python_module_version_lte():
    args = {'module': 'click', 'version': '6.8', 'comparison': 'lte'}
    check = PythonModuleCheck(**args)
    check.check()

    args['version'] = '0.0.0'

    check = PythonModuleCheck(**args)
    with pytest.raises(CheckFailure):
        check.check()


@pytest.mark.agnostic
def test_python_module_version_gt():
    args = {'module': 'click', 'version': '6.7', 'comparison': 'gt'}
    check = PythonModuleCheck(**args)
    with pytest.raises(CheckFailure):
        check.check()

    args['version'] = '6.0'

    check = PythonModuleCheck(**args)
    check.check()


@pytest.mark.agnostic
def test_python_module_version_lt():
    args = {'module': 'click', 'version': '6.8', 'comparison': 'lt'}
    check = PythonModuleCheck(**args)
    check.check()

    args['version'] = '6.7'

    check = PythonModuleCheck(**args)
    with pytest.raises(CheckFailure):
        check.check()


@pytest.mark.agnostic
def test_python_module_version_min_ver():
    args = {'module': 'click', 'min_version': '6.0', 'version': '6.9'}
    check = PythonModuleCheck(**args)
    check.check()

    args['version'] = '0.0.0'

    check = PythonModuleCheck(**args)
    with pytest.raises(CheckFailure):
        check.check()


@pytest.mark.agnostic
def test_python_module_version_invalid_args_min_ver():
    args = {'module': 'click', 'min_version': '6.7', 'comparison': 'eq'}
    with pytest.raises(CheckFailure):
        check = PythonModuleCheck(**args)


@pytest.mark.agnostic
def test_python_module_installed_only():
    args = {'module': 'click'}
    check = PythonModuleCheck(**args)
    check.check()

    args['module'] = 'not_installed'
    check = PythonModuleCheck(**args)
    with pytest.raises(CheckFailure):
        check.check()
