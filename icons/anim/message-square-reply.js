import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M22 17"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
