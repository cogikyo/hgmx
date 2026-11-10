import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const whale = el.querySelectorAll('path[d^="M6 9.7"], path[d="M15 15h.01"]');
  animate(whale, { translateY: [0, 1, -0.5, 0] }, { duration: τmax });

  const spout = el.querySelectorAll('path[d^="M18"]');
  animate(spout, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
