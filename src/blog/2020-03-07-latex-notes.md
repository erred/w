---
description: who forced us to use this thing
title: latex notes
---

### latex

oneliner (`$!` == last arg of last command):

```bash
(export f=main && pdflatex $f && bibtex $f  && pdflatex $f  && pdflatex $f)>/dev/null
```

#### text

```latex
\textbf{bold text}
\textit{italics}
\texttt{some inline code}
\underline{underlined}

\emph{context sensitive emphasis}

\footnote{some footnote here}

\caption{a caption}

\label{some:labelname}
\ref{some:labelname}
```

#### list

```latex
\begin{itemize}
  \item unordered 1
  \item unordered 2
\end{itemize}

\begin{enumerate}
  \item ordered 1
  \item ordered 2
\end{enumerate}
```

#### placement

- h: here
- t: top
- b: bottom
- !: force
- H: exactly here

##### images

```latex
\usepackage{graphicx}
\graphicspath{ {./images/} }

\begin{figure}[h]
  \centering
  \label{some:name}
  \caption{This is a caption}
  \includegraphics[width=0.4\textwidth]{name-of-image}
\end{figure}
```

##### table

after _tabular_:

- **|**: vertical line
- **l**: left
- **c**: center
- **r**: right

between lines:

- **\hline**: horizontal line

```latex
\begin{center}
\begin{tabular}{ |c|c|c| }
 \hline
 cell1 & cell2 & cell3 \\
 cell4 & cell5 & cell6 \\
 cell7 & cell8 & cell9 \\
 \hline
\end{tabular}
\end{center}
```

##### code

```latex
\usepackage{listings}
\usepackage{framed}

\lstdefinestyle{mystyle}{
  frame=single,
  basicstyle=\ttfamily\small
  breakatwhitespace=false,
  breaklines=true,
  captionpos=b,
  keepspaces=true,
  numbers=left,
  numbersep=5pt,
  showspaces=false,
  showstringspaces=false,
  showtabs=false,
  tabsize=2
}

% include from a file
\lstinputlisting{file.name}

% code block
\begin{lstlisting}[caption=some caption]
some
large
block
of
code
\end{lstlisting}
```

#### multifile

```latex
\usepackage{import}

\newpage
\import{tasks/}{task-4.tex}
```

#### bibliography

##### tex

```latex
\usepackage{biblatex}
\addbibresource{references.bib}

\cite{somereference}


\printbibliography
```

##### bib

reference: [bibliography management](https://en.wikibooks.org/wiki/LaTeX/Bibliography_Management)

```
@article{an-article,
  author  = {Author 1 and Author 2},
  title   = {Some Title},
  year    = {2020},
  journaltitle = {Some Journal},
}

@report{a-technical-report,
  author  = {Author 1},
  title   = {Some Title},
  year    = {2020},
  institution = {Some Institution},
}

@manual{a-technical-manual,
  title   = {Some Title},
  year    = {2020},
}

@online{an-online-resource,
  author  = {Author 1},
  title   = {Some Title},
  year    = {2020},
  url     = {some://url},
}

@misc{anything-else,
  author  = {Author 1},
  title   = {Some Title},
  year    = {2020},
}
```

