import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const mug = el.querySelectorAll('path[d^="M3 8"], path[d^="M17 8"]');
  animate(mug, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const string = el.querySelector('path[d^="M9 8"]');
  animate(string, { scaleY: [1, 1.5, 0.83, 1] }, { duration: τmid, delay: τmin });

  const bag = el.querySelector('path[d^="M11 16"]');
  animate(bag, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });
}
