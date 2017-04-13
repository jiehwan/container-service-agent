Name:		container-service-agent
Summary:	
Version:	0.0.1
Release:	0
Group:		System/Configuration
License:	Apache-2.0
BuildArch:	
Source0:	%{name}-%{version}.tar.gz
Source1:	%{name}.manifest
Source2:	%{name}.service

%description

%prep
%setup -q -n %{name}-%{version}
cp %{SOURCE1} ./%{name}.manifest
cp %{SOURCE2} ./%{name}.service

%build

%install

%post

%files

