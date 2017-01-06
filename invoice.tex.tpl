\documentclass{bellerophon-a}

\usepackage{tabularx}
\usepackage[table]{xcolor}

\shorttitle{Invoice}

\begin{document}

\centerline{\Huge\bf Invoice for Services Rendered}

\bigskip
\begin{minipage}[t]{3in}
Prepared by:\\
\hspace*{1em}\begin{minipage}{3in-1em}
My Company LLC\\
1234 Market St\\
Philadelphia, PA, 19111
\end{minipage}
\end{minipage}
\hfill
%
%\bigskip
\begin{minipage}[t]{3in}
Prepared for:\\
\hspace*{1em}\begin{minipage}{3in-1em}
Foo Bar\\
1234 Market St\\
Philadelphia, PA 19111
\end{minipage}
\end{minipage}

\bigskip
\hfill%\begin{minipage}[t]{3in}
\begin{tabular}[t]{ll}
{\bf Invoice \#} & FOOBAR-00\\
{\bf Invoice Period} & {{ time .Since }}--{{ time .Until }}\\
{\bf Invoice Date} & {{ time .Today }}\\
\end{tabular}
%\end{minipage}

%%----------------------------------------------------------------------
%%----------------------------------------------------------------------
\section{Invoice}

This is an invoice to Foo Bar from My Company LLC
for work performed under contract in support of Baz.

\bigskip
\noindent\begin{tabularx}{\linewidth}{|X|c|}
\hline
{\bf Item} & {\bf Qty}\\
\hline

{{ range .Entries -}}
  {{ texEscape .Title }} & {{ texDuration .TimeRounded -}}\\
{{ end }}
\hline
\hfill {\bf Total Hours} & {{ texDuration .TotalTimeRounded }}\\
\hfill {\bf Rate}        & {{ texCash .Rate }}/hr\\
\hline
\hline
\hfill {\bf Total Due} & {\bf {{ texCash .TotalDue}}}\\
\hline
\end{tabularx}

\bigskip\noindent
All values in United States dollars.

\bigskip\noindent
Checks may be made out to My Company LLC and delivered to:

\smallskip\noindent\hspace{2em}\begin{minipage}{3in}
My Company LLC\\
1234 Market St\\
Philadelphia, PA, 19111
\end{minipage}

%%----------------------------------------------------------------------
%%----------------------------------------------------------------------
\section{Contact}

Contractual as well as technical questions on this work may be addressed to:
\begin{itemize}
\item John Smith\\
      \url{smith@example.com}\\
      555-555-5555
\end{itemize}

%%----------------------------------------------------------------------
%%----------------------------------------------------------------------
\end{document}
