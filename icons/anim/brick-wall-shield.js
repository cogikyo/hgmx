import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wall = el.querySelector("g");
  animate(wall, { scaleY: [0, 1.1, 1] }, { duration: τmid });

  const shield = el.querySelector('path[d^="M22 17.5"]');
  animate(shield, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
