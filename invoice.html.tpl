<html>
<head>
    <style type="text/css">
        div.items {
            width: 100%;
            margin: auto;
        }

        .items tr, .items td {
            border-top: 1px solid lightgray;
        }

        .items table {
            border-collapse: collapse;
        }

        td.cash {
            display: flex;
        }

        span.currency {
            width: 100%;
            display: flex;
            padding-left: 1em;
        }

        span.amount {
            display: flex;
        }

        tr.summary {
            border: none;
        }

        .summary td {
            border: none;
        }
    </style>
</head>
<body style="font-family: Roboto, Helvetica, sans-serif">
<div style="width: 90%; margin: auto;">
    <h1 style="text-align: center">Invoice for Services Rendered</h1>
    <br>
    <table style="width: 100%">
        <tr style="font-size: smaller; font-weight: bold">
            <td>Prepared by:</td>
            <td>Prepared for:</td>
        </tr>
        <tr style="font-size: large">
            <td>
                My Company LLC<br>
                1234 Market St<br>
                Philadelphia, PA, 19111<br>
            </td>
            <td>Foo Bar<br>
                1234 Market St<br>
                Philadelphia, PA 19111<br>
            </td>
        </tr>
    </table>
    <br>
    <br>
    <table>
        <tr>
            <td style="font-size: smaller; font-weight: bold; width: 10em;">Invoice #</td>
            <td>FOOBAR-00</td>
        </tr>
        <tr>
            <td style="font-size: smaller; font-weight: bold; width: 10em;">Invoice Period</td>
            <td>2020/12/01 -- 2020/12/31</td>
        </tr>
        <tr>
            <td style="font-size: smaller; font-weight: bold; width: 10em;">Invoice Date</td>
            <td>2020/12/30</td>
        </tr>
    </table>

    <br>
    <br>

    <div class="items">
        <table style="width: 100%">
            <tr style="font-size: large; font-weight: bold; border-top: none;">
                <th>Description</th>
                <th style="width: 10em; max-width: 10em; text-align: center">Hours</th>
            </tr>
{{ range .Entries -}}
            <tr>
                <td>{{ .Title }}</td>
                <td style="text-align: right;">{{ duration .TimeRounded -}}</td>
            </tr>
{{ end }}
            <tr class="summary" style="border-top: 2px solid black">
                <td style="text-align: right; font-weight: bold; display: flex">
<span style="width: 50%;text-align: left"><div style="font-size: small">All currency amounts in United States dollars. Payment due net 30 days</div>
</span><span style="width: 50%">Total Hours</span></td>
                <td style="text-align: right;">{{ duration .TotalTimeRounded -}}</td>
            </tr>
            <tr class="summary">
                <td style="text-align: right; font-weight: bold">Rate (Hourly)</td>
                <td class="cash"><span class="currency">$</span><span class="amount">{{ cash .Rate -}}</span></td>
            </tr>
            <tr class="summary">
                <td style="text-align: right; font-weight: bold">Total Due</td>
                <td class="cash"><span class="currency">$</span><span class="amount">{{ cash .TotalDue -}}</span></td>
            </tr>
        </table>
    </div>

    <br>
    <div style="font-size: small">Checks may be made out to My Company LLC and delivered to:</div>
    <br>
    <div style="font-size: small; padding-left: 50px">
        My Company LLC<br>
        1234 Market St<br>
        Philadelphia, PA, 19111
    </div>
    <br>
    <div style="font-size: small">Contractual as well as technical questions on this work may be addressed to:
        <div style="font-size: small; padding-left: 50px">
            John Smith<br>
            smith@example.com<br>
            555-555-5555<br>
        </div>
    </div>
</div>
</body>
</html>