import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const chuck = el.querySelector("g");
  animate(chuck, { rotateX: [0, 360] }, { duration: τmax });

  const drill = el.querySelectorAll("path");
  animate(drill, { translateX: [0, 1.5, 0] }, { duration: τmid });
}
