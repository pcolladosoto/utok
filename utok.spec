Name:		utok
Version:	1.0
Release:	1
Summary:	Micro OpenID Connect client
BuildArch:	x86_64

URL: https://github.com/pcolladosoto/utok

License:	GPLv3

# Longer description on what the package is/does
%description
A micro client capable of generating access tokens by leveraging the
OpenID Connect Device Authorization Flow. It's been primarily designed
to work wit the Indigo IAM service but it should also be compatible
with other standards-compliant IAM implementations.

# Time to copy the binary file!
%install
# Delete the previos build root
rm -rf %{buildroot}

# Create the necessary directories
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_mandir}/man1

# And install the necessary files
install -m 0775 %{_sourcedir}/%{name}         %{buildroot}%{_bindir}/%{name}
install -m 0664 %{_sourcedir}/%{name}.1.gz    %{buildroot}%{_mandir}/man1/%{name}.1.gz

# Files provided by the package. Check https://docs.fedoraproject.org/en-US/packaging-guidelines/#_manpages too!
%files
%{_bindir}/%{name}
%{_mandir}/man1/%{name}.1*

# Changes introducd with each version
%changelog
* Mon Jan 22 2024 Pablo Collado Soto <pcolladosoto@gmx.com>
- First RPM-packaged version
