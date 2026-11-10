import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const cap = el.querySelector("g");
  animate(cap, { translateY: [0, -3, 0], rotate: [0, -15, 6, 0] }, { duration: τmax });

  const badge = el.querySelector('path[d^="M9 7"]');
  animate(badge, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
