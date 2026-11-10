import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const maker = el.querySelector('path[d^="M14 22V12"]');
  animate(maker, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const drip = el.querySelector('path[d="M7 10v2"]');
  animate(drip, { translateY: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
