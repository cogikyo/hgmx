import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const hand = el.querySelector("g");
  animate(hand, { rotate: [0, -15, 5, 0], translateY: [0, -2, 0] }, { duration: τmax });
  const outline = el.querySelector('path[d^="M15 5.88"]');
  animate(outline, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
