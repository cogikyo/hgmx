import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const shield = el.querySelector("path");
  animate(shield, { translateY: [0, -3, 1, 0] }, { duration: τmax });
  animate(shield, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
