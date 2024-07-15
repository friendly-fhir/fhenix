Feature: Download FHIR Package

Rule: Downloading FHIR packages from simplifier

  Background:
    Given a timeout of 180s
    And the registry is 'https://packages.simplifier.net'
    And the cache is empty

  Scenario: Download a Leaf FHIR Package

    Downloading a single leaf package should only download the one package and
    nothing else.

    Given the command 'fhenix download hl7.fhir.r4.core 4.0.1'
    When the command is executed
    Then the FHIR cache contains packages:
      | Package               | Version |
      | hl7.fhir.r4.core      | 4.0.1   |

  Scenario: Download a FHIR Package with Dependencies

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

  Scenario: Download a FHIR Package while excluding dependencies

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
