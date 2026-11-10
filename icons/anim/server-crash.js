import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const racks = el.querySelectorAll('path[d^="M6 10"], path[d^="M6 14"]');
  animate(racks, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const bolt = el.querySelector('path[d^="m13 6"]');
  animate(bolt, { translateX: [1.5, -0.5, 0], translateY: [-2, 0.5, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
