import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const door = el.querySelector("g");
  animate(door, { scaleX: [1, 0.8, 1] }, { duration: τmid });

  const shackle = el.querySelector('path[d^="M20 16"]');
  animate(shackle, { translateY: [0, 1.5, 0] }, { duration: τmid, delay: τmin });
}
