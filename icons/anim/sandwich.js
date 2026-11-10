import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bread = el.querySelectorAll('rect, path:not([d^="m6.67"])');
  animate(bread, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cheese = el.querySelector('path[d^="m6.67"]');
  animate(cheese, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
