import { animate, τmax } from "../motion.js";

export default function (el) {
  const elephant = el.querySelector("g");
  animate(elephant, { rotate: [0, -6, 4, -2, 0] }, { duration: τmax });
}
