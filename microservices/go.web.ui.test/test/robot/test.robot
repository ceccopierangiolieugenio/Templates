*** Settings ***
Documentation     Simple example using SeleniumLibrary.
Library           SeleniumLibrary

*** Variables ***
${BROWSER}  %{BROWSER}
${URL}      %{URL}

*** Test Cases ***
Key Press Event
    Open Browser To Index Page
    Input Email          test@test.com
    Input Card Number    12345
    Input CVV            1234
    Submit Credentials
    [Teardown]    Close Browser

*** Keywords ***
Open Browser To Index Page
    Open Browser       ${URL}    ${BROWSER}
    Title Should Be    Bootstrap 101 Template

Input Email
    [Arguments]    ${mail}
    Press Keys    inputEmail    a
    Input Text    inputEmail    ${mail}
    Press Keys    inputEmail    CTRL+C

Input Card Number
    [Arguments]    ${cn}
    Press Keys    inputCardNumber    CTRL+V
    Input Text    inputCardNumber    ${cn}

Input CVV
    [Arguments]    ${cvv}
    Input Text    inputCVV    ${cvv}

Submit Credentials
    Click Button    Submit

Welcome Page Should Be Open
    Title Should Be    Welcome Page
