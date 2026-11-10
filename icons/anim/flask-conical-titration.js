import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const flask = el.querySelector("g");
  animate(flask, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });

  const drop = el.querySelector('path[d^="M12 4"]');
  animate(drop, { translateY: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
