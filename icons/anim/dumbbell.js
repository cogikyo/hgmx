import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const dumbbell = el.querySelector("g");
  animate(dumbbell, { translateY: [0, -3, 0] }, { duration: τmid });
  animate(dumbbell, { rotate: [0, -8, 0] }, { duration: τmin, delay: τmid });
}
