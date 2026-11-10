import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bodies = el.querySelectorAll('path[d^="M21 22"], path[d^="M3 12"]');
  animate(bodies, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const arcs = el.querySelectorAll('path[d^="M20.453"], path[d^="M3.889"]');
  animate(arcs, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const heads = el.querySelectorAll("circle");
  animate(heads, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
