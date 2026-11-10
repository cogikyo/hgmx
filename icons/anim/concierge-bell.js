import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const button = el.querySelector("g g");
  animate(button, { translateY: [0, 1.5, 0] }, { duration: τmin });

  const bell = el.querySelector("g");
  animate(bell, { rotate: [0, -6, 4, 0] }, { duration: τmid, delay: τmin });
}
