import { animate, τmax } from "../motion.js";

export default function (el) {
  const ball = el.querySelector('path[d^="M15.7"]');
  animate(ball, { translateY: [5, -1, 0], rotate: [-180, 0], opacity: [0, 1] }, { duration: τmax });
}
