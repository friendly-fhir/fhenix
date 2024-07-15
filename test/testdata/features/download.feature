Feature: Download Command

The `download` command is used to explicitly download FHIR packages from a
registry for later use. The contents of the downloaded package SHALL be
cached inside a local cache directory (known as the 'FHIR Cache') for later use.

By default, downloading a package will also download all of its dependencies
unless the `--exclude-dependencies` flag is provided. If a package is already
in the cache, it will not be downloaded again unless the `--force` flag is
provided.

Rule: Downloaded FHIR packages are stored in the cache

  Downloading a FHIR package from any registry SHALL store the contents into
  the FHIR cache so that future access does not require network activity or
  expensive processing.

  Background:
    Given a timeout of 3m
    And the registry is 'https://packages.simplifier.net'
    And the cache is empty

  Scenario: FHIR Packages are stored in the FHIR Cache

    Downloading a single package that has no dependencies SHALL only download
    that one package and nothing else.

    Given the command 'fhenix download hl7.fhir.r4.core 4.0.1'
    When the command is executed
    Then the FHIR cache contains packages:
      | Package               | Version |
      | hl7.fhir.r4.core      | 4.0.1   |

Rule: Downloading FHIR packages downloads dependencies by default

  The default behavior of downloading FHIR packages SHALL include all of its
  dependencies exactly once -- unless the user manually specifies the
  `--exclude-dependencies` flag. The downloaded contents SHALL be stored into
  the FHIR Cache.

  Background:
    Given a timeout of 3m
    And the registry is 'https://packages.simplifier.net'
    And the cache is empty

  Scenario: Package contains dependencies

    Downloading a package that does have listed dependencies SHALL download
    itself, as well as all of its dependencies -- and store it into the
    FHIR Cache.

    Given the command 'fhenix download hl7.fhir.us.core 6.1.0'
    When the command is executed
    Then the FHIR cache contains packages:
      | Package                      | Version |
      | hl7.fhir.r4.core             | 4.0.1   |
      | hl7.fhir.r4.examples         | 4.0.1   |
      | hl7.fhir.us.core             | 6.1.0   |
      | hl7.fhir.uv.bulkdata         | 2.0.0   |
      | hl7.fhir.uv.extensions.r4    | 1.0.0   |
      | hl7.fhir.uv.sdc              | 3.0.0   |
      | hl7.fhir.uv.smart-app-launch | 2.1.0   |
      | hl7.terminology.r4           | 5.0.0   |
      | ihe.formatcode.fhir          | 1.1.0   |
      | us.cdc.phinvads              | 0.12.0  |
      | us.nlm.vsac                  | 0.11.0  |

  Scenario: User specifies --exclude-dependencies

    Downloading a package that has dependencies while providing the
    --exclude-dependencies flag SHALL only download the package itself and store
    it into the FHIR Cache.

    Given the command 'fhenix download hl7.fhir.us.core 6.1.0 --exclude-dependencies'
    When the command is executed
    Then the FHIR cache contains packages:
      | Package                      | Version |
      | hl7.fhir.us.core             | 6.1.0   |
    And the FHIR cache does not contain packages:
      | Package                      | Version |
      | hl7.fhir.r4.core             | 4.0.1   |
      | hl7.fhir.r4.examples         | 4.0.1   |
      | hl7.fhir.uv.bulkdata         | 2.0.0   |
      | hl7.fhir.uv.extensions.r4    | 1.0.0   |
      | hl7.fhir.uv.sdc              | 3.0.0   |
      | hl7.fhir.uv.smart-app-launch | 2.1.0   |
      | hl7.terminology.r4           | 5.0.0   |
      | ihe.formatcode.fhir          | 1.1.0   |
      | us.cdc.phinvads              | 0.12.0  |
      | us.nlm.vsac                  | 0.11.0  |

Rule: Download failure issues error

  If a package cannot be downloaded from the registry, whether by user-error
  from an invalid package name, or from a network connection, the download
  command SHALL return a non-zero status-code and provide an error message.

  Background:
    Given a timeout of 180s
    And the registry is 'https://packages.simplifier.net'
    And the cache is empty

  Scenario: Package does not exist

    Downloading a package that does not exist in the registry SHALL return a
    non-zero status-code.

    Given the command 'fhenix download hl7.fhir.r4.core 4.0.0'
    When the command is executed
    Then the program exits with non-0 status-code
    And stderr contains 'error:'

