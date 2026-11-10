import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cone = el.querySelector('path[d^="m7 11"]');
  animate(cone, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const scoop = el.querySelectorAll('path[d^="M17 7"]');
  animate(scoop, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
