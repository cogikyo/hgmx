import { animate, τmax } from "../motion.js";

export default function (el) {
  const [long, short] = el.querySelectorAll("g");
  animate(long, { translateY: [-1.5, 0.5, 0] }, { duration: τmax });
  animate(short, { translateY: [3, -0.5, 0] }, { duration: τmax });
}
