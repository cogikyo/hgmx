import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const paper = el.querySelector('path[d^="M4 3"]');
  animate(paper, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glyph = el.querySelector('path[d^="M10 17"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bars = el.querySelectorAll('path[d="M8 13h5"], path[d="M8 17h7"]');
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
