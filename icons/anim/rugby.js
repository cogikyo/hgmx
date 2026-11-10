import { animate, τmax } from "../motion.js";

export default function (el) {
  const ball = el.querySelector('path[d^="M15.7"]');
  animate(ball, { translateY: [8, -1, 0], rotate: [0, 360], opacity: [0, 1, 1] }, { duration: τmax });
}
