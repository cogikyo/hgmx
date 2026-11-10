import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glass = el.querySelector('path[d^="M20 14.54"]');
  animate(glass, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  animate(glass, { scaleY: [1, 1.15, 0.95, 1] }, { duration: τmid, delay: τmin });
  const flake = el.querySelectorAll('path:not([d^="M20 14.54"])');
  animate(flake, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
