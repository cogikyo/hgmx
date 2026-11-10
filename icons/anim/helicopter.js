import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const rotor = el.querySelector('path[d="M6 3h16"]');
  animate(rotor, { scaleX: [1, 0.2, 1, 0.2, 1] }, { duration: τmax });

  const helicopter = el.querySelector("g");
  animate(helicopter, { translateY: [0, -1.5, 0] }, { duration: τmid, delay: τmin });
}
