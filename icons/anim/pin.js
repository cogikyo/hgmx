import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="M9"]');
  animate(head, { translateY: [-6, 1, 0] }, { duration: τmax });

  const needle = el.querySelector('path[d="M12 17v5"]');
  animate(needle, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
