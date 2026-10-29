import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const outline = el.querySelector('path[d^="M17 3"]');
  animate(outline, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const facets = [el.querySelector('path[d^="M10.5"]'), el.querySelector('path[d="M2 9h20"]')];
  animate(facets, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const gem = el.querySelector("g");
  animate(gem, { scale: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
