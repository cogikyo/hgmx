import { animate, τmax } from "../motion.js";

export default function (el) {
  const [long, short] = el.querySelectorAll("g");
  animate(long, { translateX: [-1.5, 0.5, 0] }, { duration: τmax });
  animate(short, { translateX: [3, -0.5, 0] }, { duration: τmax });
}
