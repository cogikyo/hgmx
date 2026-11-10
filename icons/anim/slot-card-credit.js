import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const reader = el.querySelector("g");
  animate(reader, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });

  const stripe = el.querySelector('path[d="M11 9v10"]');
  animate(stripe, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
