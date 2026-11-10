import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const hands = el.querySelector('path[d="M12 10v2l1 1"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmax });
  const face = el.querySelector("rect");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
