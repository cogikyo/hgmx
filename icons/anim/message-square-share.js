import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M12 3"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [-4, 2, 0], translateY: [4, -2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
