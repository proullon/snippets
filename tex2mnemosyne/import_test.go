package main

import (
	"strings"
	"testing"
)

const testTex = `
\documentclass[11pt,a4paper]{article} % compile with pdfLaTeX

\usepackage{mathtools, esint, eucal}
\usepackage[utf8]{inputenc} % Accents and special characters
\usepackage[T1]{fontenc}


\begin{document}

\title{YO}
\author{Pierre Roullon}
\date{2022/11/28}
\maketitle

\pagebreak

\section{1}

\subsection{1}

\subsubsection{1}

\begin{exercise}[TD 1 Ex 3]
\begin{question}
    (a) Soient $A$ et $B$ deux ensembles finis de même cardinal, $|A| = |B|$, et soit \\
    $f: A \rightarrow B$ une application. Montrer que les propriétés suivantes sont équivalentes: \\
     (i) $f$ est injective \\
     (ii) $f$ est surjective \\
     (iii) $f$ est une bijection \\
     \par
     (b) Donner un exemple d'application injective $f: \N \rightarrow \N$ qui ne soit pas surjective, et d'une application
     surjective $g: \N \rightarrow \N$ qui ne soit pas injective. \\
     En déduire que $\N$ n'est pas fini.\\
\end{question}
\begin{solution}
    (a) Il faut montrer que si $|A| = |B|$ alors injectivité et surjectivité sont équivalentes. \\
    Supposons que $f$ est injective mais pas surjective, alors il existe $b \in B \setminus f(a)$. En particulier,
    f est aussi un injection de A dans $B \setminus {b}$. \\
    D'où $|A| \leq |B \setminus {b}| = |B| - 1$ par le principe d'injection, ce qui est une contradiction. \\
    \par
    pour la réciproque, supposons que $f$ est surjective mais pas injective. Alors il existe $a,a' \in A$ tel que $f(a) = f(a')$. \\
    En particulier, l'application $i \mapsto f(i)$ est toujours une surjection si on considère sa restriction à $A \setminus {a}$.
    D'où $|A| - 1| = |A \setminus {a}| \leq |B|$, ce qui est une contradiction. \\

    Comme f est surjective et injective, alors elle est bijective. \\
    \par
		\newcommand\test{haha}
		\test{}
    (b) $n \mapsto n + 1$ est une injection mais pas une surjection, et $n \mapsto \frac{n}{2}$ si $n$ pair et $n \mapsto 0$ sinon est une surjection non injective. \\
    Puisqu'il existe des applications non bijective de $\N$ dans $\N$, alors $\N$ n'est pas fini.
\end{solution}
\end{exercise}

\subsection{2}

\subsubsection{1}

\begin{theo}[the name of the theorem]
    Bla bli blo
\end{theo}

\subsubsection{2}

\begin{lem}[the name of the lemma]
    Haha hihi hoho
\end{lem}

\end{document}
`

func TestImportTex(t *testing.T) {

	cards, err := fromTex(testTex)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(cards) != 3 {
		t.Errorf("expected 3 cards, got %d", len(cards))
	}

	var fTheo, fEx, fLem bool
	for _, c := range cards {
		t.Logf("Tag:%s\nFront:%+s\nBack:%s\n", c.Tag, c.Front, c.Back)
		if c.Front == "<latex>Théorème: the name of the theorem</latex>" {
			fTheo = true
			if c.Back == "" {
				t.Logf("Theorem back empty")
			}
			if c.Tag != "1::2::1" {
				t.Errorf("Exercise wrong tag")
			}
		}
		if c.Front == "<latex>Lemme: the name of the lemma</latex>" {
			fLem = true
			if c.Back == "" {
				t.Logf("Lemma back empty")
			}
			if c.Tag != "1::2::2" {
				t.Errorf("Exercise wrong tag")
			}
		}
		if strings.HasPrefix(c.Front, "<latex>Exercice: ") {
			fEx = true
			if c.Back == "" {
				t.Errorf("Exercise back empty")
			}
			if c.Tag != "1::1::1" {
				t.Errorf("Exercise wrong tag")
			}
		}
	}

	if !fTheo {
		t.Errorf("Theorem not found")
	}
	if !fEx {
		t.Errorf("Exercise not found")
	}
	if !fLem {
		t.Errorf("Lemma not found")
	}

}
