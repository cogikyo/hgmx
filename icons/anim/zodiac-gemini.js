import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const arcs = el.querySelectorAll('path[d^="M20 3"], path[d^="M4 21"]');
  animate(arcs, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const pillars = [el.querySelector('path[d^="M8 4.525"]'), el.querySelector('path[d^="M16 4.525"]')];
  animate(pillars, { scaleY: [0, 1] }, { duration: τmid, delay: stagger(0.12, { startDelay: τmin }) });
}
