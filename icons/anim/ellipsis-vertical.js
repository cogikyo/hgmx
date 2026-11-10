import { animate, stagger } from "../motion.js";

export default function (el) {
  const dots = ["5", "12", "19"].map((y) => el.querySelector(`circle[cy="${y}"]`));
  animate(dots, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
