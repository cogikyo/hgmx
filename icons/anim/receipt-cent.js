import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const paper = el.querySelector('path[d^="M4 3"]');
  animate(paper, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glyph = el.querySelector('path[d^="M14.828"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bar = el.querySelector('path[d="M12 7v10"]');
  animate(bar, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
