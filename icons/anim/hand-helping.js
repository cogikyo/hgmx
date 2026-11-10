import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const hand = el.querySelectorAll("path");
  animate(hand, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
  animate(hand, { translateY: [0, -2, 0] }, { duration: τmin, delay: τmid });
}
