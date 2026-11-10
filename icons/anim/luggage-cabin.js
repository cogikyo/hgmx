import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const luggage = el.querySelector("g");
  animate(luggage, { rotate: [0, -8, 3, 0] }, { duration: τmax });

  const poles = el.querySelector("g g");
  animate(poles, { scaleY: [1, 1.3, 1] }, { duration: τmid, delay: τmin });

  const grip = el.querySelector('path[d="M8 2h8"]');
  animate(grip, { translateY: [0, -1.5, 0] }, { duration: τmid, delay: τmin });
}
