# Name of your application
APPLICATION = ninenano

# Run tests of native only
BOARD = native

# Change this to 0 show compiler invocation lines
QUIET ?= 1

# Path to the RIOT submodule
RIOTBASE ?= $(CURDIR)/../../../vendor/RIOT

# CFLAGS for debugging purposes
CFLAGS += -DDEVELHELP -ggdb

# Modules required for using sock_tcp
USEMODULE += ipv6_addr
USEMODULE += lwip_ipv6_autoconfig
USEMODULE += inet_csum
USEMODULE += lwip_ethernet lwip_netdev
USEMODULE += lwip_sock_tcp
USEMODULE += netdev_default

# Module required for unit test functions
USEMODULE += embunit

# Modules required for using the 9p component
USEMODULE += ninenano_9p

export NINEBASE := $(CURDIR)/../../..
include $(NINEBASE)/Makefile.riot
